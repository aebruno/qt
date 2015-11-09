#include "qstyleoptionbutton.h"
#include <QUrl>
#include <QModelIndex>
#include <QStyle>
#include <QStyleOption>
#include <QString>
#include <QVariant>
#include <QStyleOptionButton>
#include "_cgo_export.h"

class MyQStyleOptionButton: public QStyleOptionButton {
public:
};

void* QStyleOptionButton_NewQStyleOptionButton(){
	return new QStyleOptionButton();
}

void* QStyleOptionButton_NewQStyleOptionButton2(void* other){
	return new QStyleOptionButton(*static_cast<QStyleOptionButton*>(other));
}
