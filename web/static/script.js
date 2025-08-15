document.addEventListener('DOMContentLoaded', () => {
    const socket = new WebSocket('ws://' + window.location.host + "/ws");
    const hostLatencyMap = new Map();

    const hostsContainer = document.getElementById('hosts-container');
    let hostBoxTemplate = '';

    fetch('/templates/host-box-template.html')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.text();
        })
        .then(template => {
            hostBoxTemplate = template;
        })
        .catch(error => console.error('Error loading template:', error));

    socket.onopen = () => {
        console.log('Websocket connection established');
    };

    socket.onmessage = (event) => {
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
    }

    function updateDashboard(metrics) {
        if (!hostBoxTemplate) {
            console.warn('Template not loaded yet, skipping render.');
            return;
        }

        hostsContainer.innerHTML = ''; // Clears old data
        var barOptions = {
                chart: {
                    type: 'bar',
                    height: 350
                },
                plotOptions: {
                    bar: {
                        borderRadius: 4,
                        borderRadiusApplication: 'end',
                        horizontal: true
                    }
                },
                series: [{
                    data: Array.from(hostLatencyMap.values(), (arr) => (parseInt(arr[arr.length - 1]) / 100000))
                }],
                xaxis: {
                    categories: [...hostLatencyMap.keys()]
                }
            }

        var barChart = new ApexCharts(document.querySelector("#bar-chart"), barOptions);
        barChart.render();

        for (const host in metrics) {
            const stats = metrics[host];
            if (!hostLatencyMap.has(host)) {
                hostLatencyMap.set(host,[])
            }
            hostLatencyMap.get(host).push([stats.AvgRtt, stats.MinRtt, stats.MaxRtt])

            const isOffline = stats.PacketsRecv === 0 && stats.PacketsSent > 0;
            const statusClass = isOffline ? 'offline' : 'online';
            const statusText = isOffline ? 'Offline' : 'Online';

            const renderedHtml = hostBoxTemplate
                .replaceAll('{{host}}', host)
                .replaceAll('{{statusClass}}', statusClass)
                .replaceAll('{{statusText}}', statusText)
                .replaceAll('{{packetsSent}}', stats.PacketsSent)
                .replaceAll('{{packetsRecv}}', stats.PacketsRecv)
                .replaceAll('{{packetLoss}}', stats.PacketLoss.toFixed(1))
                .replaceAll('{{avgLatency}}', stats.AvgRtt.toFixed(2))
                .replaceAll('{{minLatency}}', stats.MinRtt.toFixed(2))
                .replaceAll('{{maxLatency}}', stats.MaxRtt.toFixed(2))
                .replaceAll('id="line-chart-{{hostName}}"', `id="line-chart-${host.replaceAll('.', '-')}"`);

            var lineChartOptions = {
                chart: {
                    type: 'line'
                },
                series: [{
                    name: "Average Latency",
                    data: hostLatencyMap.get(host).map(x => x[0])
                },
                {
                    name: "Min Latency",
                    data: hostLatencyMap.get(host).map(x => x[1])
                },
                {
                    name: "Max Latency",
                    data: hostLatencyMap.get(host).map(x => x[2])
                }],
                xaxis: {
                    labels: {
                        show: false
                    }
                }
            };

            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = renderedHtml;
            const hostBoxElement = tempDiv.firstElementChild;
            hostsContainer.appendChild(hostBoxElement);

            var lineChart = new ApexCharts(document.querySelector(`#line-chart-${host.replaceAll('.', '-')}`), lineChartOptions);
            lineChart.render();
        }
    }
})