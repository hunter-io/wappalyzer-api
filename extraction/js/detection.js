function getDetectedApps() {
  applications = [];

  for (app in wappalyzer.detected[Object.keys(wappalyzer.detected)[0]]) {
    applications.push(app);
  }

  return applications;
}
