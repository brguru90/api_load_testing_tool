const chart_option = {
    responsive: true,
    interaction: {
        mode: 'index',
        intersect: false,
    },
    stacked: false,
    plugins: {
        title: {
            align: "start",
            display: true,
            text: 'API Response Time',
            padding: {
                bottom: 40
            }
        },
        legend: {
            align: "start",
            position: "top",
            labels: {
                padding: 20
            }
        },
        tooltip: {
            callbacks: {
                label: function (context, ...r) {
                    let label = context.dataset.label || "";
                    let val = Number(context.parsed.y) || 0;
                    if (val < 1) {
                        return  label + ": " + (val * 1000) + " ms"
                    }
                    return label + ": " + val + " sec"

                }
            }
        }
    },
    scales: {
        y: {
            type: 'linear',
            display: true,
            position: 'left',
            title: {
                display: true,
                text: 'Time'
            },
            beginAtZero: true,
            ticks: {
                callback: function (label, index, labels) {
                    label = Number(label)
                    if (label < 1) {
                        return (label * 1000) + " ms"
                    }
                    return label + " sec"
                }
            }
        },
        x: {
            title: {
                display: true,
                text: 'Iterations',
                align: "start"
            },
            ticks: {
                maxRotation: 45,
                minRotation: 45,
            }
        }
    },
};

export {
    chart_option
}