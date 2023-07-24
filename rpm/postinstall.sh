if [ $1 -eq 1 ] && [ -x /usr/bin/systemctl ]; then
    # Initial installation
    systemctl enable --now etop || :
fi
