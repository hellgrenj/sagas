<script>
    import { messages } from "../stores/messages";
    import { onDestroy, onMount } from "svelte";
    import {
        Chart,
        LinearScale,
        BarController,
        CategoryScale,
        BarElement,
    } from "chart.js";

    let unsubscribe = null;
    let myChart = null;
    onMount(() => {
        Chart.register(LinearScale, BarController, CategoryScale, BarElement);
        unsubscribe = messages.subscribe((value) => {
            const labels = [];
            const data = [];
            for (const v of value) {
                addLabelIfNotExist(labels, v.Name);
            }
            labels.sort();
            for (const v of value) {
                incrementPerLabel(labels, data, v.Name);
            }
            if (myChart) {
                myChart.destroy();
            }
            RenderChart(labels, data);
        });
    });
    onDestroy(unsubscribe);
    function addLabelIfNotExist(labels, label) {
        if (!labels.includes(label)) {
            labels.push(label);
        }
    }
    function incrementPerLabel(labels, data, label) {
        let index = labels.indexOf(label);
        if (index > -1) {
            if (data[index]) {
                data[index] = data[index] + 1;
            } else {
                data[index] = 1;
            }
        }
    }
    function RenderChart(labels, data) {
        const ctx = document.getElementById("myChart").getContext("2d");
        myChart = new Chart(ctx, {
            type: "bar",
            data: {
                labels: labels,
                datasets: [
                    {
                        data: data,
                        backgroundColor: [
                            "rgba(255, 99, 132, 0.2)",
                            "rgba(54, 162, 235, 0.2)",
                            "rgba(255, 206, 86, 0.2)",
                            "rgba(75, 192, 192, 0.2)",
                            "rgba(153, 102, 255, 0.2)",
                            "rgba(255, 159, 64, 0.2)",
                        ],
                        borderColor: [
                            "rgba(255, 99, 132, 1)",
                            "rgba(54, 162, 235, 1)",
                            "rgba(255, 206, 86, 1)",
                            "rgba(75, 192, 192, 1)",
                            "rgba(153, 102, 255, 1)",
                            "rgba(255, 159, 64, 1)",
                        ],
                        borderWidth: 1,
                    },
                ],
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: true,
                    },
                },
            },
        });
    }
</script>

<main>
    <canvas id="myChart" width="300" height="300" />
</main>

<style>
    main {
        text-align: center;
        padding: 1em;
        margin: 0 auto;
    }
    #myChart {
         max-height: 500px; 
    }
</style>
