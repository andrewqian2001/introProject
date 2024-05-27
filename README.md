# introProject
Datadog intro project

# Send HTTP requests
Send request to add item:
curl --header "Content-Type: application/json" \
--request POST \
--data '{"title":"Learn go","description":"complete intern project"}' \
http://localhost:8090/addItem

Send request to get all items:
curl --header "Accept: application/json" \
--request GET \
http://localhost:8090/getAllItems

Send request to mark item as completed
curl --request PUT \
http://localhost:8090/completeItem?id=0

Notes:
- the comments are written for me to look back on, sorry if they are clogging up the space.
- I did not implement a datastore
- Still a WIP, having issues with logging at the moment, a log folder isn't being created in the container
