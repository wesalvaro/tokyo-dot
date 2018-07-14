(() => {

const GeoCodeController = function(geocode, latlng) {
  this.results = geocode(this.address);
  latlng(this.results).then((value) => {
    this.latlng = value;
  });
};

angular.module('geocode', [])
  .constant('mapsUrl', 'https://maps.googleapis.com/maps/api/')
  .constant('gApiKey', undefined)
  .factory('geocodeUrl', (gApiKey, mapsUrl) => {
    return (address, opt_apiKey) => {
      const apiKey = opt_apiKey || gApiKey;
      address = encodeURIComponent(address);
      return `${mapsUrl}geocode/json?key=${apiKey}&address=${address}`;
    }
  })
  .factory('geocode', ($http, geocodeUrl) => {
    return (address) => {
      return $http.get(geocodeUrl(address)).then((response) => {
        const results = response.data.results;
        if (!results.length) {
          throw response.data.status;
        }
        return results;
      }).catch((error) => {
        throw `"${address}": ${error}`;
      });
    };
  })
  .factory('geocodeMulti', ($q, geocode) => {
    return (addresses) => {
      return $q.all(addresses.map(geocode));
    };
  })
  .factory('latlng', (geocode) => {
    return (results, opt_index) => {
      return results.then((data) => data[opt_index || 0].geometry.location);
    }
  })
  .component('geocodeLatLng', {
    template: `
        {{ $ctrl.address }} Lat: {{ $ctrl.latlng.lat }} Lon: {{ $ctrl.latlng.lng }}
    `,
    controller: GeoCodeController,
    bindings: {
      'address': '<'
    }
  });

})();
