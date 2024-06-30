if [ $1 -eq 1 ] && [ -x /usr/bin/systemctl ]; then
    # Initial installation
    systemctl enable --now etop || :
fi

if [ $1 -eq 2 ] && [ -x /usr/bin/systemctl ]; then
    # Package upgrade
    systemctl restart etop || :
fi
