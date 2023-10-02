# ScreenRecordAPI
* Postman documentation available here [documenation](https://documenter.getpostman.com/view/29925722/2s9YJc1NWX)

## Short (quick description)  
The temporary base url is at https://screenrecordapi-production.up.railway.app

1. To start recording, send a POST request to https://screenrecordapi-production.up.railway.app/record/start
you will get an UUID in the response JSON access it with .data.uuid

2. The binary chunks should be sent to POST https://screenrecordapi-production.up.railway.app/record/upload/:uuid

3. Stop the record with POST https://screenrecordapi-production.up.railway.app/record/stop/:uuid
