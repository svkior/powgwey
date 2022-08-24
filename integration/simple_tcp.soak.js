import { check, group, sleep } from 'k6';
import tcp from 'k6/x/puretcp';

export const options = {
    stages: [
        { duration: '2m', target: 400 }, // ramp up to 400 users
        { duration: '3h56m', target: 400 }, // stay at 400 for ~4 hours
        { duration: '2m', target: 0 }, // scale down. (optional)
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