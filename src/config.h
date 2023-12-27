#ifndef CYGNUDGE_CONFIG_H
#define CYGNUDGE_CONFIG_H

#ifdef CYGNUDGE_DEBUG

#ifndef CYGNUDGE_SERVER_JSON
#define CYGNUDGE_SERVER_JSON "/home/beta-cyg/cygnudge/config/server/server.json"
#endif

#else

#ifndef CYGNUDGE_SERVER_JSON
#define CYGNUDGE_SERVER_JSON "/etc/cygnudge/server/server.json"
#endif

#endif

#ifndef CYGNUDGE_JUDGE_DELAY
#define CYGNUDGE_JUDGE_DELAY (int(200))
//unit: millisecond
#endif

#endif
