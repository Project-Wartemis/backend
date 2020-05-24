const ws = require('websocket').client;

const URL = 'https://localhost:8080/socket';
//const URL = 'https://api.wartemis.com/socket';

function start() {
  setupNewSocket(URL);
}

function setupNewSocket(endpoint) {
  let socket = new ws();

  socket.on('connectFailed', function(error) {
    console.log('connectFailed - ' + new Date());
    console.error(error.toString());
  });

  socket.on('connect', connection => {
    console.log('connect - ' + new Date());

    connection.on('error', error => {
      console.log('error - ' + new Date());
    });

    connection.on('close', () => {
      console.log('close - ' + new Date());
    });

    connection.on('message', handleMessage.bind(undefined, connection));
  });

  console.log(`connecting to socket @ ${endpoint}`);
  socket.connect(endpoint);
}

function sendMessage(connection, message) {
  console.log(`sending ${message.type} message`);
  connection.sendUTF(JSON.stringify(message));
}

function handleMessage(connection, message) {
  if(message.type !== 'utf8')
    console.log('Got a non-text message, ignoring');
  message = JSON.parse(message.utf8Data);

  console.log('message - ' + new Date());
  console.log(JSON.stringify(message));

  switch(message.type) {
    case 'connected': handleConnectedMessage(connection, message); break;
    case 'registered': handleRegisteredMessage(connection, message); break;
    case 'invite': handleInviteMessage(connection, message); break;
    case 'start': handleStartMessage(connection, message); break;
  }
}

function handleConnectedMessage(connection, message) {
  console.log('connected!');
  sendMessage(connection, {
    type: 'register',
    clientType: 'engine',
    name: 'Conquest'
  });
}

function handleRegisteredMessage(connection, message) {
  console.log(`Registered with id ${message.id}!`);
}

function handleInviteMessage(connection, message) {
  console.log(`invited to room ${message.room}!`);
  setupNewSocket(URL + '/' + message.room);
}

function handleStartMessage(connection, message) {
  console.log(`start!`);
  const state = generateInitial(message.players);
  sendMessage(connection, {
    type: 'state',
    turn: 0,
    state
  });
  let turn = 0;
  while(turn < 200) {
    generate(state);
    sendMessage(connection, {
      type: 'state',
      turn: turn++,
      state
    });
  }
  sendMessage(connection, {
    type: 'stop'
  });
}

start();

// game stuff

const NODE_COUNT = 20;

function generateInitial(players) {
  players = players.map(i => ({
    id: i,
    power: 1
  }));

  const nodes = [...Array(NODE_COUNT).keys()].map(i => ({
    id: i,
    name: 'node' + i,
    owner: -1,
    power: 2
  }));

  for(let i=0;i<players.length;i++) {
    nodes[nodes.length-1-i].owner = players[i].id;
  }

  const links = [...Array(NODE_COUNT).keys()].filter(i => i).map(i => ({
    source: i,
    target: Math.floor(Math.random() * i)
  }));

  return {
    players,
    nodes,
    links,
    events: {
      deploys: [],
      moves: []
    }
  };
}

function generate(state) {
  // deploys
  state.events.deploys = [];
  state.events.moves = [];
  for(const player of state.players) {
    player.power = state.nodes.filter(n => n.owner === player.id).length;
    if(player.power === 0) {
      continue; // if the player is dead
    }
    const target = randomNodeOfPlayer(state.nodes, player);
    state.events.deploys.push({
      target: target.id,
      power: player.power,
    });
    target.power += player.power;
  }
  // moves
  for(const player of state.players) {
    const playerNodes = state.nodes.filter(n => n.owner === player.id);
    for(const source of playerNodes) {
      const targetId = randomLinkedNodeId(source, state.links);
      const target = state.nodes.find(n => n.id === targetId);
      const power = Math.floor(source.power / 2);
      source.power -= power;
      target.power += power * (source.owner === target.owner ? 1 : -1);
      if(target.power === 0) {
        target.owner = -1;
      }
      if(target.power < 0) {
        target.owner = source.owner;
        target.power *= -1;
      }
      state.events.moves.push({
        source: source.id,
        target: target.id,
        power,
      });
    }
  }
}

function randomNodeOfPlayer(nodes, player) {
  const playerNodes = nodes.filter(n => n.owner === player.id);
  return playerNodes[Math.floor(Math.random() * playerNodes.length)];
}

function randomLinkedNodeId(node, links) {
  let neighbours = links.filter(l => l.source === node.id || l.target === node.id)
    .map(l => l.source === node.id ? l.target : l.source);
  return neighbours[Math.floor(Math.random() * neighbours.length)];
}
