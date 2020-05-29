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
  console.log('connected!');
  sendMessage(connection, {
    type: 'register',
    clientType: 'bot',
    name: 'Robbot',
    game: 'Planet Wars',
  });
}

function handleRegisteredMessage(connection, message) {
  console.log(`Registered!`);
}

function handleStateMessage(connection, message) {
  if(!message.move) {
    return;
  }
  const source = message.state.planets.find(p => p.player === 1);
  const others = message.state.planets.filter(p => p.player !== 1);
  const target = others[Math.floor(Math.random()*others.length)];
  const moves = [];
  if(source && target) {
    moves[0] = {
      source: source.id,
      target: target.id,
      ships: source.ships,
    };
  }
  sendMessage(connection, {
    type: 'action',
    game: message.game,
    key: message.key,
    action: {
      moves: moves
    }
  });
}

function handleErrorMessage(connection, message) {
  console.log(JSON.stringify(message));
}

start();
