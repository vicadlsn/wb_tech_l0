async function getOrder() {
    const orderID = document.getElementById("order_uid").value;
    const orderInfoDiv = document.getElementById("order_info");
    const errorDiv = document.getElementById("order_error");

    errorDiv.innerHTML = "";
    orderInfoDiv.style.display = "none";

    if (!orderID) {
        errorDiv.innerText = "Введите номер заказа!";
        return;
    }

    try {
        const response = await fetch(`/order/${orderID}`);
        if (!response.ok) {
            errorDiv.innerText =
                response.status === 404 ? "Заказ не найден" : "Ошибка сервера";
            return;
        }
        const data = await response.json();
        orderInfoDiv.style.display = "block";

        document.getElementById("order_info_details").innerHTML = `
            <h2>Заказ №${data.order_uid}</h2>
            <p>Трек-номер: ${data.track_number}</p>
            <p>Дата оформления: ${data.date_created}</p>
            <p>Служба доставки: ${data.delivery_service}</p>
        `;

        const d = data.delivery;
        document.getElementById("order_info_delivery").innerHTML = `
            <h3>Доставка</h3>
            <p>Имя: ${d.name}</p>
            <p>Телефон: ${d.phone}</p>
            <p>Email: ${d.email}</p>
            <p>Адрес: ${d.region}, ${d.city}, ${d.address}, ${d.zip}</p>
        `;

        const p = data.payment;
        document.getElementById("order_info_payment").innerHTML = `
            <h3>Оплата</h3>
            <p>Товары: ${p.goods_total} ${p.currency}</p>
            <p>Доставка: ${p.delivery_cost} ${p.currency}</p>
            <p>Пошлина: ${p.custom_fee} ${p.currency}</p>
            <p>Банк: ${p.bank}</p>
            <p>Платёжный сервис: ${p.provider}</p>
            <p><strong>Сумма: ${p.amount} ${p.currency}</strong></p>
        `;

        let itemsHTML = `
                    <h3>Товары</h3>
                    <table>
                        <thead>
                            <tr>
                                <th>№</th>
                                <th>Название</th>
                                <th>Бренд</th>
                                <th>Цена</th>
                                <th>Скидка</th>
                                <th>Итог</th>
                                <th>Размер</th>
                            </tr>
                        </thead>
                        <tbody>
                `;
        if (data.items && data.items.length !== 0) {
            data.items.forEach((item, i) => {
                itemsHTML += `
                <tr>
                    <td>${i + 1}</td>
                    <td>${item.name}</td>
                    <td>${item.brand}</td>
                    <td>${item.price}</td>
                    <td>${item.sale}%</td>
                    <td>${item.total_price}</td>
                    <td>${item.size}</td>
                </tr>
            `;
            });
        }

        itemsHTML += "</tbody></table>";
        document.getElementById("order_info_items").innerHTML = itemsHTML;
    } catch (err) {
        console.error(err);
        errorDiv.innerText = "Ошибка при загрузке данных";
    }
}

document.addEventListener("DOMContentLoaded", () => {
    document
        .getElementById("get_order_info_button")
        .addEventListener("click", getOrder);
});
