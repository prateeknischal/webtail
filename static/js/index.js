angular.module("webtail").controller("mainController", mainController);

mainController.$inject = ["$rootScope", "$scope", "$mdSidenav"]

function mainController($rootScope, $scope, $mdSidenav) {
  var vm = this;

  // $scope.toggleLeft = vm.toggleSideNav('left')
  // $scope.toggleRight = vm.toggleSideNav('right')

  vm.toggleSideNav = function toggleSideNav() {
    $mdSidenav('left').toggle()
  }

  vm.init = function init() {
    console.log("In the main controller")
  }

  vm.init();
}
