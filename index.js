var child_process = require('child_process');
var byline = require('./byline');

exports.handler = function(event, context) {
  var proc = child_process.spawn('./main', [JSON.stringify(event)], { stdio: [process.stdin, 'pipe', 'pipe'] });

  proc.stdout.on('data', function(line){
    var msg = JSON.parse(line);
    context.succeed(msg);
  })

  proc.on('exit', function(code){
    console.error('exit blabla: %s', code)
    process.exit(1)
  })
}