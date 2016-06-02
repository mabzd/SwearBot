lock '3.4.0'

set :application, "swearbot"
set :repo_url, "git@github.com:mabzd/SwearBot.git"

set :linked_files, %w(bin/token.txt bin/log.txt bin/config.json bin/settings.json bin/swears.txt)
set :keep_releases, 5
set :deploy_to, "/var/go/swearbot"

namespace :app do
  task :compile do
    on roles(:app) do
      execute "cd #{current_path} && mkdir -p bin"
      execute "cd #{current_path} && cp -u swears.txt bin/swears.txt"
      execute "cd #{current_path} && cp -u config-rename.json bin/config.json"
      execute "cd #{current_path} && GOPATH=$HOME/go go build -o ./bin/swbot.exe main.go"
    end
  end

  %i(start stop).each do |command|
    task command do
      on roles(:app) do
        execute "/etc/init.d/swbot #{command}"
      end
    end
  end
end

after "deploy", "app:compile"
after "deploy", "app:stop"
after "deploy", "app:start"
