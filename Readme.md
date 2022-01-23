# Trash App Backend.
This is the backend service for the [React Native Trash App](https://github.com/frankbo/trash-app)
## Running
The backend runs in AWS Lambda. The service can run locally with the `-local` parameter. `go run main.go -local`
## Parameters
`locationId` for example 1746.18
`streetId` for example 1746.34.1

Those ids come from [The Trash Calendar Bad Berleburg](https://www.bad-berleburg.de/Alte-Struktur/Leben-in-Bad-Berleburg/Wohnen-Umwelt/Abfallkalender/index.php?ort=&strasse=1746.34.1&vtyp=2&vMo=01&bMo=12&vJ=2022&call=sfm&La=1&css=&bn=&Barriere=&sNavID=1746.8&mNavID=1746.8&ffmod=abf&ffsm=1)