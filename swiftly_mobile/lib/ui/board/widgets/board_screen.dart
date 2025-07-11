import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/board/widgets/cart_widget.dart';

class BoardScreen extends StatelessWidget {
  const BoardScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      backgroundColor: Color.fromARGB(255, 9, 30, 114),
      body: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Padding(
          padding: EdgeInsets.symmetric(horizontal: 10, vertical: 20),
          child: Row(
            children: [
              IntrinsicWidth(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('To do', style: TextStyle(color: Colors.white, fontSize: 18),),
                    SizedBox(height: 10),
                    CartWidget(categoryWidget: CategoryWidget(name: 'аналитика'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.high,),
                    SizedBox(height: 10),
                    CartWidget(categoryWidget: CategoryWidget(name: 'программирование'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.medium,),
                    SizedBox(height: 10),
                    CartWidget(categoryWidget: CategoryWidget(name: 'аналитика'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.low,),
                  ],
                ),
              ),
              SizedBox(width: 20),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('In progress', style: TextStyle(color: Colors.white, fontSize: 18),),
                  SizedBox(height: 10),
                  CartWidget(categoryWidget: CategoryWidget(name: 'дизайн'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.medium,),
                  SizedBox(height: 10),
                  CartWidget(categoryWidget: CategoryWidget(name: 'аналитика'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.high,),
                  SizedBox(height: 10),
                  CartWidget(categoryWidget: CategoryWidget(name: 'программирование'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.high,),
                  SizedBox(height: 10),
                  CartWidget(categoryWidget: CategoryWidget(name: 'программирование'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.medium,),
                ],
              ),
              SizedBox(width: 20),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('Completed', style: TextStyle(color: Colors.white, fontSize: 18),),
                  SizedBox(height: 10),
                  CartWidget(categoryWidget: CategoryWidget(name: 'дизайн'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.low,),
                  SizedBox(height: 10),
                  CartWidget(categoryWidget: CategoryWidget(name: 'аналитика'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.high,),
                ],
              ),
              SizedBox(width: 20),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('In review', style: TextStyle(color: Colors.white, fontSize: 18),),
                  SizedBox(height: 10),
                  CartWidget(categoryWidget: CategoryWidget(name: 'программирование'), name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug', priority: Priority.low,),
                ],
              ),
            ],
          ),
        ),
      )
    );
  }
}

//CartWidget(name: 'Lorem ipsum', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit', data: '5 aug'),