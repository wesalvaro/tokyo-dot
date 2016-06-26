(function loadScripts() {
  const scripts = [
    'https://cdnjs.cloudflare.com/ajax/libs/d3/3.4.11/d3.min.js',
    'https://cpettitt.github.io/project/graphlib-dot/v0.4.10/graphlib-dot.min.js',
  ];
  scripts.forEach(function(script) {
    const e = document.createElement('script');
    e.src = script;
    document.body.appendChild(e);
  });
}());

let geocode = function(apiKey) {
  const URL_GEOCODE = `https://maps.googleapis.com/maps/api/geocode/json?key=${apiKey}&address=`;
  const notFound = {};
  const found = {};
  const locations = [];

  trainGraph.then((graph) => {
    graph.eachNode(function(id, node) {
      if (!node.label || node.label.match(/^\w/)) return;
      locations.push({id: node.id, name: node.label});
    });

    locations.forEach(function(location) {
      d3.json(URL_GEOCODE + '' + location.name + '駅', function(e, loc) {
        if (e) {
          notFound[location.id] = JSON.parse(e.responseText).error_message;
          return;
        }
        if (!loc.results.length) {
          notFound[location.id] = loc.status;
          return;
        }
        let latlon = loc.results[0].geometry.location;
        found[location.id] = latlon;
      });
    });
  });

  return {
    locations: locations,
    found: found,
    notFound: notFound,
  };
};

let getGeocoder = function(apiKey) {
  const result = geocode(apiKey);
  return function() {
    console.log(`${Object.keys(result.found).length}/${result.locations.length}`);
    console.log(result.notFound);
    return `const GEOCODE = ${JSON.stringify(result.found)};`;
  };
}
