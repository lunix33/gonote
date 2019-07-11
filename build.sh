# To successfully build the project you need to have the following components.
# GCC (64bit compatible) -- This is to be able to compile the SQLite module.
#	For Windows you'll need to use mingw-w64
#	Be sure gcc is in the system $PATH

# This script can be run on Windows with Git Bash terminal.

cdir=$(cd $(dirname $0); pwd)

# Required environment variables.
export CGO_ENABLED=1
export GOPATH=$cdir

# Frontend.
echo "**** Building frontend..."
cd $cdir/ngNote/
echo "*** Installing dependancies ..."
npm install
echo "** Building scripts..."
npx ng build --prod

# Backend.
echo "**** Building backend..."
cd $cdir/src/gonote/
echo "*** Installing dependancies..."
go get
echo "** Compile files..."
go install

echo "**** Collecting build files..."
# Cleanup
if [ -d "$cdir/build" ]; then
	echo "*** Cleaning up old files..."
	mv $cdir/build/notes.db $cdir
	rm -fvr $cdir/build/
fi
mkdir -p $cdir/build/public/
mv $cdir/notes.db $cdir/build/

echo "** Copy executable ..."
cp $cdir/bin/gonote.exe $cdir/build/
echo "* Copy frontend files..."
cp $cdir/ngNote/dist/ngNote/* $cdir/build/public/

echo "Build complete."
echo "Files are available in: $cdir/build"
