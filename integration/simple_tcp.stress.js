import { check, group, sleep } from 'k6';
import tcp from 'k6/x/puretcp';

export const options = {
    stages: [
        { duration: '2m', target: 100 }, // below normal load
        { duration: '5m', target: 100 },
        { duration: '2m', target: 200 }, // normal load
        { duration: '5m', target: 200 },
        { duration: '2m', target: 300 }, // around the breaking point
        { duration: '5m', target: 300 },
        { duration: '2m', target: 400 }, // beyond the breaking point
        { duration: '5m', target: 400 },
        { duration: '10m', target: 0 }, // scale down. Recovery stage.
    ],
};

export default function () {
    const conn = tcp.connect('server:8000')
    //tcp.write(conn, 'some data\n');
    const reply = tcp.getQuote(conn)
    //console.log('reply',reply)
    check(reply, {
        'Non null answer': (resp) => resp.length > 0,
    });
}