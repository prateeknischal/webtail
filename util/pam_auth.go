package util

// #cgo LDFLAGS: -lpam
/*

// Code taken from : https://stackoverflow.com/questions/10910193/how-to-authenticate-username-password-using-pam-w-o-root-privileges
#include <stdio.h>
#include <security/pam_appl.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>

struct pam_response *reply;

// //function used to get user input
int function_conversation(int num_msg, const struct pam_message **msg, struct pam_response **resp, void *appdata_ptr){
    *resp = reply;
        return PAM_SUCCESS;
}

int authenticate_system(const char *username, const char *password)
{
    const struct pam_conv local_conversation = { function_conversation, NULL };
    pam_handle_t *local_auth_handle = NULL; // this gets set by pam_start

    int retval;
    retval = pam_start("su", username, &local_conversation, &local_auth_handle);

    if (retval != PAM_SUCCESS){
        return 0;
    }

    reply = (struct pam_response *)malloc(sizeof(struct pam_response));

    reply[0].resp = strdup(password);
    reply[0].resp_retcode = 0;
    retval = pam_authenticate(local_auth_handle, 0);

    if (retval != PAM_SUCCESS){
        return 0;
    }
    retval = pam_end(local_auth_handle, retval);

    if (retval != PAM_SUCCESS){
        return 0;
    }
    return 1;
}

*/
import "C"

// PamAauthenticate - authenticates a user based on local PAM based
// authentication
func PamAauthenticate(username, password string) int {
	uname := C.CString(username)
	pass := C.CString(password)
	res := C.int(C.authenticate_system(uname, pass))
	return int(res)
}
