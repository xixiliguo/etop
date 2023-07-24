if [ $1 -eq 0 ] && [ -x /usr/bin/systemctl ]; then
    # Package removal, not upgrade
    systemctl disable --now etop || :
fi
