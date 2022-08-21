import { check, group, sleep } from 'k6';
import tcp from 'k6/x/puretcp';

export const options = {
    stages: [
        { duration: '1m', target: 100 }, // simulate ramp-up of traffic from 1 to 100 users over 5 minutes.
        { duration: '3m', target: 100 }, // stay at 100 users for 10 minutes
        { duration: '1m', target: 0 }, // ramp-down to 0 users
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