Badbank
=======

Protokoll och implementation av Badbank TCP protokollet. Skolarbete i kursen programmeringsparadigmer.

Protokollet delas upp i två portar. Query för använder trafik och Update för klientuppdateringar. Båda är Big-Endian.

| Del    | Port | Typ            |
| -----  | ---- | -------------- |
| Query  | 1337 | 10 byte binary |
| Update | 1338 | JSON           |

## Query

Query delen hanterar allt som har med användare att göra. Att logga in, felmeddelanden, kontoförändringar osv. Query är ett binär protokoll med ett enda format uppdelat i fyra områden.

| Område    | Längd (bitar) | Tolkning     |
| --------- | :-----------: | ------------ |
| opcode    | 2             | unsigned int |
| special   | 14            | *olika       |
| storklump | 64            | signed int   |

### Opcode
Operationskoden talar om vad meddelandets typ är. Det finns 4 olika opcodes.

| Opcodes | namn   | binärkod |
| :-----: | ------ | -------- |
| 0       | iam*   | 00       |
| 1       | login  | 01       |
| 2       | change | 10       |
| 3       | info   | 11       |
* skickas över Update porten. 

#### Iam
Iam skickas för att tala om vilket språk klienten vill ha och trigga en nerladdning av den senaste klientdatan så som språk och välkomstmeddelande. 
Iam är speciell på två sätt. Den skickas över Update porten 1338 och har en underlig special. Iam's special är 2*7 ascii, alltså komprimerad ascii, där extended biten på ascii bokstäverna är borttagna. Iam's special innehåller vilket språk klienten ska vara på i typen av sv, en, uk, osv.

| Format  | opcode | special | storklump |
| ------- | ------ | ------- | --------- |
| områden | iam    | språk   | 0         |
| exempel | 00     | "sv"    | 0         |

#### Login
Login är det första meddelandet klienten måste skicka på Query porten. Pinkod och kontonummer är argumenten.

| Format  | opcode | special | storklump   |
| ------- | ------ | ------- | ----------- |
| områden | login  | pinkod  | kontonummer |
| exempel | 01     | 4444    | 3141592654  |

#### Change
Efter att användaren har loggat in med Login kan den skicka Change frågor för att ändra på användarens konto. Change kan både lägga till och tabort gotyckligt stora summor från kontot i minsta valören av användarens hemlandsvaluta, så som ören eller cent. Change måste skickas med en tvåsiffrig säkerhetskod i special för att transaktionen ska gå igenom. Säkerhetskoden är alla ojämna tal mellan 0 - 100.

| Format  | opcode | special      | storklump |
| ------- | ------ | ------------ | --------- |
| områden | change | säkerhetskod | mängd     |
| exempel | 10     | 01           | -299      |

#### Info
Info används för att skicka informationsmeddelanden så som felmeddelanden eller saldon.

| Specialfältet | namn                  |
| ------------- | --------------------- |
| 2             | ok                    | 
| 3             | bad login             |
| 4             | error                 |
| 5             | internt fel           |
| 6             | bad verification code |

| Format  | opcode | special       | storklump |
| ------- | ------ | ------------  | --------- |
| områden | info   | meddelandetyp | saldo     |
| exempel | 11     | 2             | 300       |


### Exempeltrafik i Query

| Klient    | Riktning | Server |
| --------- | :------: | ------ |
| Login     | -->      |        |
|           | <--      | info   |
| Change    | -->      |        |
|           | <--      | info   |
| Change    | -->      |        |
|           | <--      | info   |
| TCP CLOSE |          |        |


## Update
Uppdateringar från server kan komma när som helst och triggas av olika saker. Skickar användaren en IAM till server eller om server bara har en ny version som den just nu håller på att trycka ut.

Uppdateringar kommer i form av [JSON](http://www.json.org/) encodad data.

| Fält               | Typ                    |
| ------------------ | ---------------------- |
| Välkomstmeddelande | 80 tecken sträng utf-8 |
| TODO...            |                        |

### Exempeltrafik i Update

| Klient    | Riktning | Server |
| --------- | :------: | ------ | 
| Iam       | -->      |        |
|           | <--      | update |
| ...       | ...      | ...    |
|           | <--      | update |
| TCP CLOSE |          |        |