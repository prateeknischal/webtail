angular.module("webtail").controller("loginController", loginController);

loginController.$inject = ["$rootScope", "$scope"]

function loginController($rootScope, $scope) {
  var vm = this;

  $scope.errDialog = false
  function getUrlVars(){
    var vars = [], hash;
    var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
    for(var i = 0; i < hashes.length; i++)
    {
        hash = hashes[i].split('=');
        vars.push(hash[0]);
        vars[hash[0]] = hash[1];
    }
    return vars;
  }
  vm.init = function init() {
    console.log("In the login controller")
    var query = getUrlVars();
    console.log(query);
    if (query["err"] === "invalid") {
      $scope.errDialog = true;
    }
  }

  vm.init();
}
