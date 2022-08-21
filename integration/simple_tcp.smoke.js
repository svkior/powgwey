import { check, group, sleep } from 'k6';
import tcp from 'k6/x/puretcp';

export const options = {
    stages: [
        { duration: '10s', target: 1 },
    ]
};

export default function () {
    const conn = tcp.connect('example8:8972')
    //tcp.write(conn, 'some data\n');
    const reply = tcp.read(conn)
    //console.log('reply',reply)
    check(reply, {
        'Non null answer': (resp) => resp.length > 0,
    });
}