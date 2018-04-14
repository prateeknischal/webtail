angular.module('webtail', ["ngMaterial"])
.config(["$qProvider", function($qProvider){
  $qProvider.errorOnUnhandledRejections(false);
}]);

angular.module('webtail', ['ngMaterial'])
.config(function($mdThemingProvider) {
  $mdThemingProvider.theme('default')
    .primaryPalette('orange')
    .accentPalette('amber');
  $mdThemingProvider.setDefaultTheme('default');
});
