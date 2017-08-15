(() => {

const sigmaGraph = () => {
  return {
    scope: {
      nodes: '<',
      edges: '<',
      type: '<',
    },
    link: (s, e) => {
      s.sigma = new sigma();
      s.$watch('type', (type) => {
        type = type || 'canvas';
        if (s.sigma.renderers.main) s.sigma.killRenderer('main');
        s.sigma.addRenderer({
          id: 'main',
          type: type,
          container: e[0],
        });
        s.sigma.render();
      });
      s.$watchGroup(['nodes', 'edges'], (values) => {
        const [nodes, edges] = values;
        if (!nodes || !edges) return;
        s.sigma.graph.clear();
        s.sigma.graph.read({
          nodes: nodes,
          edges: edges,
        });
        s.sigma.refresh();
      });
    },
  };
};

const sigmaDotGrapher = () => {
  return (data, opt_labelParser, opt_getXY) => {
    const labelParser = opt_labelParser || ((x) => x);

    return {
      nodes: data.nodes().map(id => {
        const node = data.node(id);
        if (!node.label) return;
        const n = {
          id: id,
          type: 'station',
          size: 1,
          label: labelParser(node.label, node)
        };
        if (opt_getXY) {
          [n.x, n.y] = opt_getXY(id, node);
        }
        return n;
      }).filter(n => !!n),
      edges: data.edges().map(e => {
        const edge = data.edge(e);
        const {v: source, w: target} = e;
        return {
          id: e.name || `${source}-${target}`,
          size: 1,
          color: d3.rgb(edge.color).toString(),
          source: source,
          target: target,
          type: (edge.color == '#000000' || edge.color == '#0000ff') ? undefined : 'curve',
        };
      }),
    };
  };
};

angular.module('sigmaGraphs', [])
  .directive('sigmaGraph', sigmaGraph)
  .factory('sigmaDotGrapher', sigmaDotGrapher)
  ;

})();
