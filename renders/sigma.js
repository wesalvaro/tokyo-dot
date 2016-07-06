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
    const nodes = [];
    const edges = [];
    const labelParser = opt_labelParser || ((x) => x);

    data.eachNode(function(id, node) {
      if (!node.label) return;
      node.id = id;
      node.type = 'station';
      node.size = 1;
      node.label = labelParser(node.label, node);
      if (opt_getXY) {
        [node.x, node.y] = opt_getXY(id, node);
      }
      nodes.push(node);
    });

    data.eachEdge(function(id, source, target, edge) {
      edge.id = id;
      edge.size = 1;
      edge.color = d3.rgb(edge.color).toString();
      edge.source = source;
      edge.target = target;
      edge.type = (edge.color == '#000000' || edge.color == '#0000ff') ? undefined : 'curve';
      edges.push(edge);
    });

    return {nodes: nodes, edges: edges};
  };
};

angular.module('sigmaGraphs', [])
  .directive('sigmaGraph', sigmaGraph)
  .factory('sigmaDotGrapher', sigmaDotGrapher)
  ;

})();
