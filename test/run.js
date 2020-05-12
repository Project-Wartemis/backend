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

    register(connection);
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
    name: 'Robbot',
    key: KEY
  });
}

function sendRandomEcho(connection) {
  if(!connection.connected)
    return;
  sendMessage(connection, {
    type: 'echo',
    value: Math.round(Math.random() * 0xFFFFFF)+''
  });
  setTimeout(sendRandomEcho, 1000, connection);
}

function handleMessage(connection, message) {
  if(message.type !== 'utf8')
    console.log('Got a non-text message, ignoring');
  message = JSON.parse(message.utf8Data);

  console.log('Got a message!');
  console.log(JSON.stringify(message));

  switch(message.type) {
    case 'game': handleGameMessage(connection, message); break;
  }
}

function handleGameMessage(connection, message) {
  console.log('handled');
  setupNewSocket(URL + '/' + message.key);
}

start();
