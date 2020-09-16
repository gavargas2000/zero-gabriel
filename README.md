# zero-gabriel

To be able to run make sure you have the following packages:

go get -u github.com/gorilla/mux

go get -u github.com/gorilla/websockets



Endpoint #1:
ws://127.0.0.1:8844/websocket/SESSION_ID


Endpoint#2:
http://localhost:8844/session/SESSION_ID



I used the index.html file to test the Websocket endpoint,
added the file to the project in case you want to use it too and just change the json objects.
You can of course use it several times with the same session_id, and different elements as described
on the Assigment.

For some reason I could not make it to work with Curl via console, not sure if that was the intended method of testing,
but at least the provided index.html file allows to test, or maybe you have another way of testing it that hopefully
works.

The code contains inline comments for some of the decisions made.
