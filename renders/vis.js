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

    return {
      nodes: data.nodes().map(id => {
        const node = data.node(id);
        if (!node.label) return;
        const n = {
          id: id,
          label: labelParser(node.label),
          shape: 'box',
          group: id.match(/^[mA-Z]+/)[0],
          color: node.color,
        };
        if (opt_getXY) {
          [n.x, n.y] = opt_getXY(id);
        }
        return n;
      }).filter(n => !!n),
      edges: data.edges().map(e => {
        const edge = data.edge(e);
        const {v: from, w: to} = e;
        return {
          dashes: Boolean(edge.style == 'dashed'),
          from: from,
          to: to,
          color: edge.color,
        };
      }),
    };
  };
};

angular.module('visGraphs', [])
  .directive('visGraph', visGraph)
  .factory('visDotGrapher', visDotGrapher)
  ;

})();
