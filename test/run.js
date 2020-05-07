const ws = require('websocket').client;
const uuid = require('uuid');

let socket = new ws();

socket.on('connect', connection => {
  console.log('connected!');

  connection.on('error', error => {
    console.log('error: ' + error);
  });

  connection.on('close', () => {
    console.log('closed');
  });

  connection.on('message', handleMessage);

  register(connection)
  //sendRandomEcho(connection);
  sendGamestate(connection);
});

socket.connect('ws://localhost:8080/socket');

function sendMessage(connection, message) {
  connection.sendUTF(JSON.stringify(message));
}

function register(connection) {
  sendMessage(connection, {
    type: 'register',
    name: 'Robbot',
    key: uuid.v4()
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

let turn = 0
function sendGamestate(connection) {
  if(!connection.connected)
    return;
  sendMessage(connection, {
    type: 'gamestate',
    payload: {
      players: ['me', 'myself', 'i'],
      turn: turn++,
      stuff: {
        random: 'text'
      }
    }
  });
  setTimeout(sendGamestate, 3000, connection);
}

function handleMessage(message) {
  if(message.type !== 'utf8')
    return;
  console.log(JSON.stringify(message.utf8Data))
}
