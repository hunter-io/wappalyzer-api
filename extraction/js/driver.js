(function() {
	if ( typeof wappalyzer === 'undefined' ) {
		return;
	}

	var
		w             = wappalyzer,
		debug         = true,
		d             = window.document,
		domain        = window.top.location.host,
		url           = window.top.location.href.replace(/#.*$/, ''),
		hasOwn        = Object.prototype.hasOwnProperty,
    env           = [],
    headers       = {};

	w.driver = {
		timeout: 1000,

		/**
		 * Log messages to console
		 */
		log: function(args) {
			if ( debug && console != null && console[args.type] != null ) {
				console[args.type](args.message);
			}
		},

		/**
		 * Initialize
		 */
		init: function() {
			w.driver.getEnvironmentVars();
			w.driver.getResponseHeaders();

      w.analyze(domain, url, { html: d.documentElement.innerHTML, env: env, headers: headers });
		},

		getEnvironmentVars: function() {
			w.log('func: getEnvironmentVars');

			var i;

			for ( i in window ) {
				env.push(i);
			}
		},

		getResponseHeaders: function() {
			w.log('func: getResponseHeaders');

			var xhr = new XMLHttpRequest();

			xhr.open('GET', url, true);

			xhr.onreadystatechange = function() {
				if ( xhr.readyState === 4 && xhr.status ) {
					var headers = xhr.getAllResponseHeaders().split("\n");

					if ( headers.length > 0 && headers[0] != '' ) {
						w.log('responseHeaders: ' + xhr.getAllResponseHeaders());

						headers.forEach(function(line) {
							var name, value;

							if ( line ) {
								name  = line.substring(0, line.indexOf(': '));
								value = line.substring(line.indexOf(': ') + 2, line.length - 1);

								headers[name.toLowerCase()] = value;
							}
						});
					}
				}
			}

			xhr.send();
		},

		/**
		 * Display apps
		 */
		displayApps: function() {
      w.log('func: displayApps');

      var app, cats, apps  = [];

			for ( app in wappalyzer.detected[url] ) {
				cats = [];

				wappalyzer.apps[app].cats.forEach(function(cat) {
					cats.push(wappalyzer.categories[cat].name);
				});

				apps.push({
					name: app,
					confidence: wappalyzer.detected[url][app].confidenceTotal.toString(),
					version:    wappalyzer.detected[url][app].version,
					icon:       wappalyzer.apps[app].icon || 'default.svg',
					website:    wappalyzer.apps[app].website,
					categories: cats
				});
			}

			return apps;
		}
	};

	w.init();
})();
