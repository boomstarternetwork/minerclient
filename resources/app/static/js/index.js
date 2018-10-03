document.addEventListener('astilectron-ready', function() {

    var App = function() {
        var self = this;

        self.projects = ko.observableArray();
        self.currencies = ko.observableArray();

        self.project = ko.observable();
        self.currency = ko.observable();
        self.address = ko.observable();
        self.worker = ko.observable();

        self.mining = ko.observable();
        self.logLines = ko.observableArray();

        self.mineable = ko.pureComputed(function() {
            var project = self.project(),
                currency = self.currency(),
                address = self.address();
            return project && currency && address;
        });

        self.getProjects = function() {
            astilectron.sendMessage({"name": "getProjects"}, function(msg) {
                if (msg.name === "error") {
                    self.addLogLine({
                        "type": "err",
                        "text": "Error! Failed to get projects: "+msg.payload
                    });
                } else {
                    self.projects(msg.payload)
                }
            });
        };

        self.getCurrencies = function() {
            astilectron.sendMessage({"name": "getCurrencies"}, function(msg) {
                if (msg.name === "error") {
                    self.addLogLine({
                        "type": "err",
                        "text": "Error! Failed to get currencies: "+msg.payload
                    });
                } else {
                    self.currencies(msg.payload)
                }
            });
        };

        self.init = function() {
            app.getProjects();
            app.getCurrencies();
        };

        self.startMining = function() {
            var msg = {
                "name": "startMining",
                "payload": {
                    "projectId": self.project(),
                    "currency": self.currency(),
                    "address": self.address(),
                    "worker": self.worker()
                }
            };
            astilectron.sendMessage(msg, function(msg) {
                if (msg.name === "error") {
                    self.addLogLine({
                        "type": "err",
                        "text": "Error! Failed to start mining: "+msg.payload
                    });
                } else {
                    self.mining(true);
                }
            })
        };

        self.stopMining = function() {
            astilectron.sendMessage({"name": "stopMining"}, function(msg) {
                if (msg.name === "error") {
                    self.addLogLine({
                        "type": "err",
                        "text": "Error! Failed to stop mining: "+msg.payload
                    });
                } else {
                    self.mining(false);
                }
            })
        };

        self.addLogLine = function(line) {
            self.logLines.unshift(line);
            if (self.logLines().length > 100) {
                self.logLines().pop()
            }
        };


        astilectron.onMessage(function(msg) {
            switch (msg.name) {
                case "error":
                    self.addLogLine({"type":"err", "text":msg.payload});
                    break;
                case "logLine":
                    self.addLogLine({"type":"out", "text":msg.payload});
                    break;
                default:
                    self.addLogLine({
                        "type": "err",
                        "text": "Unknown message: "+msg.name+", payload: "+
                            JSON.stringify(msg.payload)
                    });
                    break;
            }
        });

        return self;
    };

    app = new App();

    ko.applyBindings(app);

    app.init();
});