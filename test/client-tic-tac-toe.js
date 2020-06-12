const ws = require('websocket').client;

const URL = 'https://localhost:8080/socket';
//const URL = 'https://api.wartemis.com/socket';

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

function handleMessage(connection, message) {
  if(message.type !== 'utf8')
    console.log('Got a non-text message, ignoring');
  message = JSON.parse(message.utf8Data);

  switch(message.type) {
    case 'connected': handleConnectedMessage(connection, message); break;
    case 'registered': handleRegisteredMessage(connection, message); break;
    case 'state': handleStateMessage(connection, message); break;
    case 'error': handleErrorMessage(connection, message); break;
  }
}

function handleConnectedMessage(connection, message) {
  sendMessage(connection, {
    type: 'register',
    clientType: 'bot',
    name: 'Robbot',
    game: 'Tic Tac Toe',
  });
}

function handleRegisteredMessage(connection, message) {
  console.log(`Registered with id ${message.id}`);
}

function handleStateMessage(connection, message) {
  console.log(message);
  if(!message.move) {
    return;
  }
  sendMessage(connection, {
    type: 'action',
    game: message.game,
    key: message.key,
    action: {
      position: message.state.board.indexOf(" ")
    }
  });
}

function handleErrorMessage(connection, message) {
  console.log(JSON.stringify(message));
}

start();
