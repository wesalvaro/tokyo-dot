const GRAPH_FILE = '../tokyo.dot';

const trainGraph = new Promise((resolve, reject) => {
  d3.text(GRAPH_FILE, function(e, dot) {
    if (e) reject(e);
    else resolve(graphlibDot.parse(dot));
  });
});
