(() => {
    'use strict';

    const basePath = new URL('../', document.currentScript.src).pathname.replace(/\/$/, '');
    const queryURL = (service) => `${window.location.origin}${basePath}/${service}/1/query`;
    const previewLimit = 64 * 1024;
    const selectors = ['network', 'station', 'location', 'channel'];
    const rangeTimes = ['starttime', 'endtime'];
    const epochTimes = ['startbefore', 'startafter', 'endbefore', 'endafter'];
    const geographicFields = {
        none: [],
        rectangle: ['minlatitude', 'maxlatitude', 'minlongitude', 'maxlongitude'],
        circle: ['latitude', 'longitude', 'minradius', 'maxradius']
    };
    const decimal = /^[+-]?(?:\d+(?:\.\d*)?|\.\d+)$/;
    const quoteShell = (value) => `'${value.replaceAll("'", "'\\''")}'`;
    const formatSize = (size) => {
        if (size < 1024) return `${size} B`;
        if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KiB`;
        return `${(size / (1024 * 1024)).toFixed(2)} MiB`;
    };
    const serviceReference = (service) => ({
        curl: [
            'curl --fail-with-body',
            '--header "Authorization: Bearer $FDSN_ACCESS_KEY"',
            "--header 'Content-Type: text/plain'",
            '--data-binary @selections.txt',
            `--output ${service === 'station' ? 'station.xml' : 'waveform.mseed'}`,
            quoteShell(queryURL(service))
        ].join(' \\\n  ')
    });
    const defaultParameters = (service) => {
        const end = new Date(Math.floor(Date.now() / 1000) * 1000);
        return {
            network: '*',
            station: '*',
            location: '*',
            channel: '*',
            starttime:
                service === 'dataselect'
                    ? new Date(end.getTime() - 60000).toISOString().slice(0, 19)
                    : '',
            endtime: service === 'dataselect' ? end.toISOString().slice(0, 19) : '',
            startbefore: '',
            startafter: '',
            endbefore: '',
            endafter: '',
            level: 'station',
            format: service === 'station' ? 'xml' : 'miniseed',
            includerestricted: true,
            minlatitude: '',
            maxlatitude: '',
            minlongitude: '',
            maxlongitude: '',
            latitude: '',
            longitude: '',
            minradius: '',
            maxradius: '',
            quality: 'B',
            minimumlength: '',
            longestonly: false,
            nodata: '204'
        };
    };

    const urlBuilder = (service) => ({
        service,
        ...defaultParameters(service),
        timeMode: 'range',
        geographicMode: 'none',
        copyMessage: '',
        copyError: '',
        accessKey: '',
        requesting: false,
        requestController: null,
        requestError: '',
        responseStatus: '',
        responseFailed: false,
        responseURL: '',
        responseType: '',
        responseSize: '',
        responseText: '',
        responseTruncated: false,
        responseMessage: '',
        downloadURL: '',
        downloadName: '',

        get requestDisabled() {
            return this.requesting || Boolean(this.validationError);
        },

        get requestLabel() {
            return this.requesting ? 'Sending…' : 'Send request';
        },

        get responseStatusClass() {
            return this.responseFailed ? 'app-form-status-error' : '';
        },

        get fileName() {
            if (service === 'dataselect') return 'waveform.mseed';
            return this.format === 'xml' ? 'station.xml' : 'station.txt';
        },

        get timeFields() {
            return this.timeMode === 'range' ? rangeTimes : epochTimes;
        },

        get geographicFields() {
            return geographicFields[this.geographicMode];
        },

        get url() {
            const query = new URLSearchParams();
            const names = [...selectors, ...this.timeFields];
            if (service === 'station') {
                names.push('level', 'format', 'includerestricted', ...this.geographicFields);
            } else {
                names.push('quality', 'minimumlength', 'longestonly', 'format');
            }
            names.push('nodata');
            for (const name of names) {
                let value = String(this[name]).trim();
                if (value === '') continue;
                if ([...rangeTimes, ...epochTimes].includes(name)) {
                    const time = new Date(`${value}Z`);
                    if (Number.isFinite(time.getTime())) value = time.toISOString();
                }
                query.set(name, value);
            }
            return `${queryURL(service)}?${query}`;
        },

        get curl() {
            return [
                'curl --fail-with-body',
                '--header "Authorization: Bearer $FDSN_ACCESS_KEY"',
                `--output ${quoteShell(this.fileName)}`,
                quoteShell(this.url)
            ].join(' \\\n  ');
        },

        get validationError() {
            for (const name of selectors) {
                const value = this[name].trim();
                if (service === 'dataselect' && !value) return `${name} is required.`;
                if (/[^\x20-\x7e]/.test(value)) {
                    return `${name} must contain only printable ASCII characters.`;
                }
            }
            for (const name of this.timeFields) {
                const value = this[name];
                if (!value && service === 'dataselect') return `${name} is required (UTC).`;
                if (value && !Number.isFinite(Date.parse(`${value}Z`))) {
                    return `${name} must be a valid UTC date and time.`;
                }
            }
            if (
                this.timeMode === 'range' &&
                this.starttime &&
                this.endtime &&
                this.starttime > this.endtime
            ) {
                return 'End time must not be earlier than start time.';
            }
            if (service === 'station') {
                if (this.level === 'response' && this.format !== 'xml') {
                    return 'Instrument responses require XML format.';
                }
                for (const name of this.geographicFields) {
                    const value = String(this[name]).trim();
                    if (!value) {
                        if (name === 'latitude' || name === 'longitude') {
                            return 'Latitude and longitude are both required for a radius search.';
                        }
                        continue;
                    }
                    const limit = name.includes('latitude') ? 90 : 180;
                    const minimum = name.includes('radius') ? 0 : -limit;
                    if (!decimal.test(value) || Number(value) < minimum || Number(value) > limit) {
                        return `${name} must be a decimal number between ${minimum} and ${limit}.`;
                    }
                }
                for (const dimension of ['latitude', 'longitude', 'radius']) {
                    const min = String(this[`min${dimension}`]).trim();
                    const max = String(this[`max${dimension}`]).trim();
                    if (
                        this.geographicFields.includes(`min${dimension}`) &&
                        min &&
                        max &&
                        Number(min) > Number(max)
                    ) {
                        return `Minimum ${dimension} must not exceed maximum ${dimension}.`;
                    }
                }
            } else {
                const minimum = String(this.minimumlength).trim();
                if (
                    minimum &&
                    (!decimal.test(minimum) ||
                        !Number.isFinite(Number(minimum)) ||
                        Number(minimum) < 0)
                ) {
                    return 'Minimum length must be a nonnegative decimal number of seconds.';
                }
            }
            if (new TextEncoder().encode(this.url).length > 2000) {
                return 'The URL is too long. Shorten the selectors or use a bulk POST request.';
            }
            return '';
        },

        syncFormat() {
            if (this.level === 'response') this.format = 'xml';
        },

        clearFeedback() {
            this.copyMessage = '';
            this.copyError = '';
        },

        reset() {
            if (this.requesting) return;
            Object.assign(this, defaultParameters(service));
            this.timeMode = 'range';
            this.geographicMode = 'none';
            this.copyMessage = '';
            this.copyError = '';
            this.requestError = '';
            this.clearResponse();
        },

        clearResponse() {
            if (this.downloadURL) URL.revokeObjectURL(this.downloadURL);
            this.downloadURL = '';
            this.downloadName = '';
            this.responseStatus = '';
            this.responseFailed = false;
            this.responseURL = '';
            this.responseType = '';
            this.responseSize = '';
            this.responseText = '';
            this.responseTruncated = false;
            this.responseMessage = '';
        },

        cancelRequest() {
            this.requestController?.abort();
        },

        destroy() {
            this.cancelRequest();
            this.clearResponse();
            this.accessKey = '';
        },

        async sendRequest() {
            if (this.requesting) return;
            this.requestError = '';
            if (this.validationError) {
                this.requestError = this.validationError;
                return;
            }
            const key = this.accessKey.trim();
            if (!key) {
                this.requestError = 'Enter an access key.';
                return;
            }

            this.clearResponse();
            this.requesting = true;
            const controller = new AbortController();
            this.requestController = controller;
            const url = this.url;
            const fileName = this.fileName;
            try {
                const response = await fetch(url, {
                    headers: { Authorization: `Bearer ${key}` },
                    mode: 'same-origin',
                    credentials: 'omit',
                    cache: 'no-store',
                    redirect: 'error',
                    signal: controller.signal
                });
                const body = await response.blob();
                const preview = service === 'station' || !response.ok;
                const text = preview ? await body.slice(0, previewLimit).text() : '';
                if (controller.signal.aborted) return;

                this.responseStatus = `HTTP ${response.status} ${response.statusText}`.trim();
                this.responseFailed = !response.ok;
                this.responseURL = url;
                this.responseType = response.headers.get('Content-Type') || 'No content type';
                this.responseSize = formatSize(body.size);
                this.responseText = text;
                this.responseTruncated = preview && body.size > previewLimit;
                this.responseMessage = response.ok && !body.size ? 'No data returned.' : '';
                if (response.ok && body.size) {
                    this.downloadURL = URL.createObjectURL(body);
                    this.downloadName = fileName;
                }
            } catch {
                this.requestError = controller.signal.aborted
                    ? 'Request canceled.'
                    : 'Unable to complete the request. Check your connection and try again.';
            } finally {
                this.requesting = false;
                this.requestController = null;
            }
        },

        async copy(value, label) {
            this.copyMessage = '';
            this.copyError = '';
            if (this.validationError) return;
            try {
                if (navigator.clipboard) {
                    await navigator.clipboard.writeText(value);
                } else {
                    const input = document.createElement('textarea');
                    input.value = value;
                    input.className = 'fdsn-copy-input';
                    input.setAttribute('aria-hidden', 'true');
                    input.tabIndex = -1;
                    document.body.append(input);
                    try {
                        input.select();
                        if (!document.execCommand('copy')) throw new Error('Copy failed');
                    } finally {
                        input.remove();
                    }
                }
                this.copyMessage = `${label} copied.`;
            } catch {
                this.copyError = 'Unable to access the clipboard. Select and copy the text below.';
            }
        },

        copyURL() {
            return this.copy(this.url, 'URL');
        },

        copyCurl() {
            return this.copy(this.curl, 'cURL command');
        }
    });

    document.addEventListener('alpine:init', () => {
        Alpine.data('fdsnHeader', () => ({
            menuOpen: false,

            get navigationClass() {
                return this.menuOpen ? 'app-navigation-open' : '';
            },

            get menuLabel() {
                return this.menuOpen ? 'Close navigation menu' : 'Open navigation menu';
            },

            get menuIcon() {
                return this.menuOpen ? 'M6 6l12 12M6 18L18 6' : 'M4 6h16M4 12h16M4 18h16';
            },

            toggleMenu() {
                this.menuOpen = !this.menuOpen;
            },

            resetMenu() {
                this.menuOpen = false;
            },

            closeMenu() {
                if (!this.menuOpen) return;
                this.resetMenu();
                this.$refs.menuToggle.focus();
            }
        }));
        Alpine.data('stationBuilder', () => urlBuilder('station'));
        Alpine.data('dataselectBuilder', () => urlBuilder('dataselect'));
        Alpine.data('stationReference', () => serviceReference('station'));
        Alpine.data('dataselectReference', () => serviceReference('dataselect'));
    });
})();
