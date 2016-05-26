var storage;
try {
  storage = JSON.parse(localStorage['lol-personal-counter']);
} catch (e) {
  storage = {};
}
var $form = document.querySelector('form');
var $regionBtn = document.querySelector('.btn.btn-default');
var $regionList = document.querySelector('.dropdown-menu');
var $spinner = document.querySelector('.spinner-area');
var $header = document.querySelector('body > .header');
var $container = document.querySelector('body > .container');
var $errorContainer = document.querySelector('.alert.alert-danger');

// analytics
var $role = document.querySelector('.role');
var $summonerName = document.querySelector('.summonerName');
var $yawhideLink = document.querySelector('p a');

var $originalErrorContainer = $errorContainer.innerHTML;

if ($form) {
  if (storage.summonerName) {
    $form.name.value = storage.summonerName;
    $regionBtn.innerText = storage.region || 'NA';
    $form.region.value = $regionBtn.innerText;
    $form.role.value = storage.role || 'Top';
    $form.save.checked = true;
    $form.enemy.focus();
    console.info('Username saved:', storage.summonerName);
  }

  $form.onsubmit = function (e) {
    e.preventDefault();
    if ($form.save.checked) {
      storage.summonerName = $form.name.value;
      storage.region = $regionBtn.innerText.trim();
      storage.role = $form.role.value;
    } else if (localStorage['lol-personal-counter']) {
      delete storage.summonerName;
      delete storage.region;
      delete storage.role;
    }
    localStorage['lol-personal-counter'] = JSON.stringify(storage);
    // turn on spinner
    $spinner.style.display = 'block';
    $form.style.display = 'none';
    $errorContainer.style.display = 'none';
    $errorContainer.innerHTML = $originalErrorContainer;
    var status;
    var now = new Date();

    fetch('/matchup', {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        SummonerName: $form.name.value,
        Enemy       : $form.enemy.value,
        Role        : $form.role.value,
        Region      : $form.region.value,
        RememberMe  : $form.save.checked,
        Referral    : document.referrer,
      })
    })
      .then(function(response) {
        status = response.status;
        return response.text()
      })
      .then(function(body) {
        var time = new Date() - now;
        if (time - now < 400) {
          return setTimeout(handle, 400 - time);
        }
        function handle() {
          if (status === 400 || status === 500|| status == 404) {
            $errorContainer.innerHTML += body;
            $errorContainer.style.display = 'block';
            $form.style.display = 'block';
            $spinner.style.display = 'none';
            return;
          }
          document.body.removeChild($container);
          document.querySelector('.response').insertAdjacentHTML('afterbegin', body);
          loadMatchupHandlers();

          ga('send', 'pageview', '/matchup');
        }
      });

    fetch('/analytics/index', {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        SummonerName: $form.name.value,
        Enemy       : $form.enemy.value,
        Role        : $form.role.value,
        RememberMe  : $form.save.checked,
        Referral    : document.referrer,
      })
    });
  }

  $regionBtn.onclick = function () {
    $regionList.style.display = !$regionList.style.display || $regionList.style.display === 'none' ? 'block' : 'none';
  }

  $regionList.onmousedown = function (e) {
    $form.region.value = e.target.innerText;
    $regionBtn.innerText = e.target.innerText;
    $regionList.style.display = 'none';
  }

  $regionList.onblur = function () {
    $regionList.style.display = 'none';
  }

  $regionBtn.onblur = function (e) {
    $regionList.style.display = 'none';
  }
}

if ($yawhideLink) {
  $yawhideLink.onclick = function (e) {
    fetch('/analytics/external', {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        Url : e.target.href,
        Page: window.location.pathname,
      })
    });
  }
}

function loadMatchupHandlers() {
  var $championGGLinks = document.querySelectorAll('table tbody tr th a');
  var $slider = document.querySelector('.slider.round');
  var $sliderInput = document.querySelector('.switch input[type="checkbox"]');
  var $switchLabel = document.querySelector('.switch-label');
  var hideInexperiencedChamps = 'Hide inexperienced champions';
  var showInexperiencedChamps = 'Show inexperienced champions';

  if (storage.showInexperienced) {
    var elems = document.querySelectorAll('#champion-level-1');
    for (var i = 0; i < elems.length; i++) {
      elems[i].style.display = 'table-row';
    }
    $sliderInput.checked = false;
    $switchLabel.innerText = hideInexperiencedChamps;
  }

  $slider.onclick = function (e) {
    var elems = document.querySelectorAll('#champion-level-1');
    console.log(elems)
    if ($switchLabel.innerText.toLowerCase() === showInexperiencedChamps.toLowerCase()) {
      for (var i = 0; i < elems.length; i++) {
        elems[i].style.display = 'table-row';
      }
      $switchLabel.innerText = hideInexperiencedChamps;
      storage.showInexperienced = true;
      localStorage['lol-personal-counter'] = JSON.stringify(storage);
    } else {
      for (var i = 0; i < elems.length; i++) {
        elems[i].style.display = 'none';
      }
      $switchLabel.innerText = showInexperiencedChamps;
      delete storage.showInexperienced;
      localStorage['lol-personal-counter'] = JSON.stringify(storage);
    }
  }

	var bgC = new RGBA(255, 255, 255, 0);

	var champList = document.getElementById("championList").getElementsByTagName("tr");

	for (var i = champList.length - 1; i > 0; i--) {
		var percent = parseInt(champList[i].children[1].innerText);
		setBG(champList[i], percent);
		if (i == 2) bgC = new RGBA(255, 245, 153, 0);
		if (i == 3) bgC = new RGBA(188, 219, 246, 0);
		if (i == 4) bgC = new RGBA(169, 239, 222, 0);
	}

  for (var i = 0; i < $championGGLinks.length; i++) {
    $championGGLinks[i].onclick = function (e) {
      fetch('/analytics/matchup', {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          SummonerName: $summonerName.id,
          Enemy       : e.target.innerText,
          Role        : $role.id,
          Click       : e.target.href,
        })
      });
    }
  }

	function RGBA(red, green, blue, alpha) {
		this.red = red;
		this.green = green;
		this.blue = blue;
		this.alpha = alpha;
		this.getCSS = function() {
			return "rgba(" + this.red + ", " + this.green + ", " + this.blue + ", " + this.alpha + ")";
		}
	}

	function setBG(elem, opac) {
		bgC.alpha = opac/2/100;
		elem.style.backgroundColor = bgC.getCSS();
	}
}
