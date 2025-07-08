import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';

class LeftPanel extends StatelessWidget {
  const LeftPanel({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 10, vertical: 20),
      child: const Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Text('Ваше облако', style: TextStyle(color: AppColors.white),),
              Spacer(),
              Icon(Icons.dehaze, color: AppColors.white,),
            ],
          ),
          SizedBox(height: 12),
          Row(
            children: [
              Icon(Icons.folder, color: AppColors.grey,),
              SizedBox(width: 5),
              Text('Все файлы', style: TextStyle(color: AppColors.grey,),),
            ],
          ),
          SizedBox(height: 8),
          Row(
            children: [
              Icon(Icons.folder, color: AppColors.grey,),
              SizedBox(width: 5),
              Text('Последние', style: TextStyle(color: AppColors.grey,),),
            ],
          ),
          SizedBox(height: 8),
          Row(
            children: [
              Icon(Icons.folder, color: AppColors.grey,),
              SizedBox(width: 5),
              Text('По проектам', style: TextStyle(color: AppColors.grey,),),
            ],
          ),
          SizedBox(height: 8),
          Row(
            children: [
              Icon(Icons.folder, color: AppColors.grey,),
              SizedBox(width: 5),
              Text('Избранное', style: TextStyle(color: AppColors.grey,),),
            ],
          ),
          SizedBox(height: 8),
          Row(
            children: [
              Icon(Icons.folder, color: AppColors.grey,),
              SizedBox(width: 5),
              Text('По задачам', style: TextStyle(color: AppColors.grey,),),
            ],
          ),
        ],
      ),
    );
  }
}