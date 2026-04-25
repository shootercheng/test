const myMap = new Map([
    ['name', 'John'],
    ['age', 30]
]);
const jsonString = JSON.stringify(Object.fromEntries(myMap));
console.log(jsonString);
