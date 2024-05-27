# introProject
Datadog intro project

# Send HTTP requests
Send request to add item:

```
curl --header "Content-Type: application/json" \
--request POST \
--data '{"title":"Learn go","description":"complete intern project"}' \
http://localhost:8090/addItem
```

Send request to get all items:
```
curl --header "Accept: application/json" \
--request GET \
http://localhost:8090/getAllItems
```
Send request to mark item as completed
```
curl --request PUT \
http://localhost:8090/completeItem?id=0
```
# Set API Key
Create a .env file and set API Key

DD_API_KEY=key_name

# Notes:
- the comments are written for me to look back on, sorry if they are clogging up the space.
- I did not implement some optional features such as datastore, dd-trace-go 
- Still a WIP

# Questions:
1. having issues with logging, a log folder isn't being created in the container

![image](https://github.com/andrewqian2001/introProject/assets/51491033/2ac0728e-25a8-472f-9206-17ebd09274f6)

   
shouldn't this line of code in main.go (line 76) create a log folder if it doesnt exist?
![image](https://github.com/andrewqian2001/introProject/assets/51491033/3c90ebe6-2f92-40ef-8fc4-d187bca5b4eb)





