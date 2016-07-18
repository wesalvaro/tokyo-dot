(() => {

const visGraph = () => {
  return {
    scope: {
      nodes: '<',
      edges: '<',
    },
    link: (s, e) => {
      const options = {
        nodes: {
          shape: 'box',
          size: 10,
          font: {
            size: 22,
            color: '#fff',
          },
          borderWidth: 2,
        },
        edges: {
          width: 2,
        },
      };
      s.network = new vis.Network(e[0]);
      s.network.setOptions(options);
      s.network.on('startStabilizing', () => {
        // Stop the drifting!
        setTimeout(() => s.network.stopSimulation(), 7000);
      });
      s.$watchGroup(['nodes', 'edges'], (values) => {
        const [nodes, edges] = values;
        if (!nodes || !edges) return;
        s.network.setData({nodes: nodes, edges: edges});
        s.network.redraw();
      });
    },
  };
};

const visDotGrapher = () => {
  return (data, opt_labelParser, opt_getXY) => {
    const labelParser = opt_labelParser || ((x) => x);
    const nodes = [];
    data.eachNode(function(id, node) {
      if (!node.label) return;
      node.id = id;
      node.label = labelParser(node.label);
      node.shape = 'box';
      node.group = id.match(/^[mA-Z]+/)[0];
      if (opt_getXY) {
        [node.x, node.y] = opt_getXY(id);
      }
      nodes.push(node);
    });

    const edges = [];
    data.eachEdge(function(id, from, to, edge) {
      edge.dashes = Boolean(edge.style == 'dashed');
      edge.from = from;
      edge.to = to;
      edges.push(edge);
    });

    return {nodes: nodes, edges: edges};
  };
};

angular.module('visGraphs', [])
  .directive('visGraph', visGraph)
  .factory('visDotGrapher', visDotGrapher)
  ;

})();
