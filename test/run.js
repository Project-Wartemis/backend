const ws = require('websocket').client;
const uuid = require('uuid');

const URL = 'https://localhost:8080/socket';
//const URL = 'https://api.wartemis.com/socket';
const KEY = uuid.v4();

function start() {
  setupNewSocket(URL);
}

function setupNewSocket(endpoint) {
  let socket = new ws();

  socket.on('connectFailed', function(error) {
    console.error('Connect Error: ' + error.toString());
  });

  socket.on('connect', connection => {
    console.log('connected!');

    connection.on('error', error => {
      console.log('error: ' + error);
    });

    connection.on('close', () => {
      console.log('closed');
    });

    connection.on('message', handleMessage.bind(undefined, connection));
  });

  console.log(`connecting to socket @ ${endpoint}`);
  socket.connect(endpoint);
}

function sendMessage(connection, message) {
  connection.sendUTF(JSON.stringify(message));
}

function register(connection) {
  sendMessage(connection, {
    type: 'register',
    clientType: 'bot',
    name: 'Robbot',
    key: KEY
  });
}

function handleMessage(connection, message) {
  if(message.type !== 'utf8')
    console.log('Got a non-text message, ignoring');
  message = JSON.parse(message.utf8Data);

  console.log('Got a message!');
  console.log(JSON.stringify(message));

  switch(message.type) {
    case 'connect': handleConnectMessage(connection, message); break;
    case 'game': handleGameMessage(connection, message); break;
  }
}

function handleConnectMessage(connection, message) {
  sendMessage(connection, {
    type: 'register',
    clientType: 'bot',
    name: 'Robbot',
    key: KEY
  });
}

function handleGameMessage(connection, message) {
  console.log('handled');
  setupNewSocket(URL + '/' + message.key);
}

start();
