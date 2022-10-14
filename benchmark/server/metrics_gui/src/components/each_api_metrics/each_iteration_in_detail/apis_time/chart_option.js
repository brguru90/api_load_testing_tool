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
            position: "bottom",
            labels: {
                padding: 20
            }
        },
    },
    scales: {
        y: {
            type: 'linear',
            display: true,
            position: 'left',
            title:{
                display: true,
                text: 'Time (ms)'
            },
            beginAtZero: true
        },
        x:{
            title:{
                display: true,
                text: 'Iterations',
                align:"start"
            }
        }
    },
};

export {
    chart_option
}