document.getElementById('calcForm').addEventListener('submit', function(event) {
    event.preventDefault();
    var input = document.getElementById('input').value;
    var result = evaluateExpression(input);
    document.getElementById('result').textContent = 'Результат: ' + result;
});