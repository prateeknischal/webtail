angular.module('webtail', ["ngMaterial"])
.config(["$qProvider", function($qProvider){
  $qProvider.errorOnUnhandledRejections(false);
}]);

angular.module('webtail', ['ngMaterial'])
.config(function($mdThemingProvider) {
  $mdThemingProvider.theme('default')
    .primaryPalette('blue')
    .accentPalette('light-blue');
});
