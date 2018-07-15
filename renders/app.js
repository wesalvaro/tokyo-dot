(() => {

const tokyoDotController = function(
    $scope, addStationColor, trainDotGraphJP, sigmaDotGrapher, visDotGrapher) {
  this.sid = 'TY';
  $scope.$watchGroup(['$ctrl.sid', '$ctrl.sNodes', '$ctrl.sEdges'], (values) => {
    const [id, nodes, edges] = values;
    if (!nodes || !edges) return;
    this.sigmaNodes = nodes.filter(n => n.id.startsWith(id));
    this.sigmaEdges = edges.filter(e => e.source.startsWith(id) && e.target.startsWith(id));
  });
  this.graph = trainDotGraphJP;
  this.graph.then((g) => {
    let {nodes, edges} = sigmaDotGrapher(g, undefined, (id) => {
      return getXY(id);
    });
    this.sNodes = nodes;
    this.sEdges = edges;
    ({nodes, edges} = visDotGrapher(g, undefined, (id) => {
      return getXY(id);
    }));
    this.vNodes = nodes;
    this.vEdges = edges;
  });
};

angular.module('app', ['sigmaGraphs', 'visGraphs', 'geocode'])
  .constant('trainDotUrl', '../tokyo.dot')
  .constant('dotParser', (graphStr) => graphlibDot.read(graphStr))
  .factory('stationGeocodes', (geocodeMulti, latlng, trainDotGraphJP) => {
    return trainDotGraphJP.then((graph) => {
      const stations = [];
      graph.nodes().forEach(id => {
        const node = graph.node(id);
        if (!node.label || node.label.match(/^\w/)) return;
        stations.push({id: id, name: node.label.replace(' ', '駅')});
      });
      return geocodeMulti(stations.map((station) => station.name)).then(
          (ll) => ll.map((l) => l[0].geometry.location)).then(
              (ll) => {
                const geocodes = {};
                stations.forEach((station, i) => {
                  geocodes[station.id] = ll[i];
                });
                return geocodes;
              });
    });
  })
  .factory('trainDotGraphJP', (trainDotGraph, addStationColor) => {
    return trainDotGraph.then((graph) => {
      graph.nodes().forEach(id => {
        const node = graph.node(id);
        const label = node.label;
        if (!label) return;
        const match = label.match(/\{(.*?)\|(.*?)\}\|\{\w+\|\w+\}/);
        if (!match) {
          console.log(label)
        }
        node.label = match[1];
        node.type = 'station';
        addStationColor(node.label, node.color);
      });
      return graph;
    });
  })
  .factory('trainDotGraph', ($http, dotParser, trainDotUrl) => {
    return $http.get(trainDotUrl).then((result) => {
      return dotParser(result.data);
    });
  })
  .factory('addStationColor', () => {
    const stationColors = {};
    sigma.canvas.nodes.station = (node, context, settings) => {
      const prefix = settings('prefix') || '';
      const size = node[prefix + 'size'];
      const x = node[prefix + 'x'];
      const y = node[prefix + 'y'];

      let start = 0;
      stationColors[node.label].forEach((color, i, arr) => {
        const end = ((i + 1) / arr.length) * 2 * Math.PI;
        context.fillStyle = color || settings('defaultNodeColor');
        context.beginPath();
        context.moveTo(x, y);
        context.arc(x, y, size, start, end);
        context.lineTo(x, y);
        context.fill();
        start = end;
      });
    };
    return (label, color) => {
      stationColors[label] = stationColors[label] || [];
      stationColors[label].push(color);
    };
  })
  .component('tokyoDot', {
    template: `
      <div ng-init="sigma=true">
        <input type="search" ng-model="$ctrl.sid">
        <input type="checkbox" ng-model="sigma">
        Actual?
      </div>
      <sigma-graph
          nodes="$ctrl.sigmaNodes"
          edges="$ctrl.sigmaEdges"
          ng-if="sigma"></sigma-graph>
      <vis-graph
          nodes="$ctrl.vNodes"
          edges="$ctrl.vEdges"
          ng-if="!sigma"></vis-graph>
    `,
    controller: tokyoDotController,
  })
  /* Enable to update Geocode data:
  .constant('gApiKey', 'API_KEY_HERE')
  .run((stationGeocodes) => {
    stationGeocodes.then((results) => {
      console.log(`const GEOCODE = ${JSON.stringify(results)}`);
    });
  })
  */
  ;

})();
