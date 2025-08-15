document.addEventListener('DOMContentLoaded', () => {
    const hostsContainer = document.getElementById('hosts-container');
    const socket = new WebSocket('ws://' + window.location.host + "/ws");
    let hostBoxTemplate = '';

    // const TIMEOUT_DURATION = 5000;
    // let timeoutID;

    // function resetTimeout() {
    //     clearTimeout(timeoutID);
    //     timeoutID = setTimeout(handleTimeout, TIMEOUT_DURATION)
    // }

    // function handleTimeout() {
    //     console.error('Websocket timeout; No data received for', TIMEOUT_DURATION, 'ms');
    // }

    fetch('/static/host-box-template.html')
        .then(response => response.text())
        .then(template => {
            hostBoxTemplate = template
        })
        .catch(error => console.error('Error loading template:', error));

    socket.onopen = () => {
        console.log('Websocket connection established');
        // resetTimeout();
    };

    socket.onmessage = (event) => {
        // resetTimeout();
        try {
            const data = JSON.parse(event.data)
            console.log('Recieved data:', data)
            updateDashboard(data);
        } catch (e) {
            console.error('Failed to parse JSON:', e);
        }
    }

    socket.onclose = () => {
        console.log('Websocket connection closed');
        // handleTimeout();
    }

    function updateDashboard(metrics) {
        hostsContainer.innerHTML = ''; // Clears old data

        for (const host in metrics) {
            const hostMetrics = metrics[host];
            const stats = hostMetrics.Stats;
            const recentLoss = hostMetrics.RecentLoss;

            let statusColor = 'green';
            let statusText = 'Online'

            if (recentLoss === 100) {
                statusColor = 'red';
                statusText = 'Timeout/Offline';
            } else if (recentLoss > 0) {
                statusColor = 'orange';
                statusText = 'Packet Loss';
            }

            // Create a new div for the host's stats
            let renderedHtml = hostBoxTemplate
                .replace('{{host}}', host)
                .replace('{{statusColor}}', statusColor)
                .replace('{{statusText}}', statusText)
                .replace('{{packetsSent}}', stats.PacketsSent)
                .replace('{{packetsRecv}}', stats.PacketsRecv)
                .replace('{{packetLoss}}', stats.PacketLoss.toFixed(1))
                .replace('{{avgLatency}}', stats.AvgRtt.toFixed(2))
                .replace('{{minLatency}}', stats.MinRtt.toFixed(2))
                .replace('{{maxLatency}}', stats.MaxRtt.toFixed(2));

            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = renderedHtml;
            hostsContainer.appendChild(tempDiv.firstElementChild);
            oldPacketLossDif = dif;
        }
    }
})