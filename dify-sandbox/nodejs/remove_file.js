const fs = require('fs');
const path = require('path');

function emptyDir(dir) {
    fs.readdirSync(dir).forEach(f => {
        const p = path.join(dir, f);
        fs.statSync(p).isDirectory() ? emptyDir(p) : fs.unlinkSync(p);
    });
}

emptyDir(process.cwd());
console.log('file removed');