lock '3.4.0'

set :application, "swearbot"
set :repo_url, "git@github.com:mabzd/SwearBot.git"

set :linked_files, %w(bin/config.json log.txt)
set :keep_releases, 5
set :deploy_to, "/var/go/swearbot"

namespace :app do
  task :compile do
    on roles(:app) do
      execute "cd #{current_path} && mkdir -p bin"
      execute "cd #{current_path} && cp swears.txt bin/swears.txt"
      execute "cd #{current_path} && GOPATH=/home/michal/go go build -o ./bin/swbot.exe main.go"
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
after "deploy", "app:start"
