import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../domain/models/label_item.dart';
import '../../../providers/label_notifier_provider.dart';
import '../../kanban/widgets/card_details_widgets/button_widget.dart';
import '../../kanban/widgets/card_details_widgets/text_field_widget.dart';
import '../themes/colors.dart';
import '../ui/custom/custom_color_picker.dart';

class LabelItemSettings extends ConsumerWidget {
  final LabelItem labelItem;
  LabelItemSettings({super.key, required this.labelItem});

  final _titleController = TextEditingController();
  late  Color currentColor;

  void changeColor(Color color) {
      currentColor = color;
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    _titleController.text = labelItem.title;
    currentColor = labelItem.color;

    return ClipRRect(
      borderRadius: BorderRadius.circular(12),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
        child: Container(
          decoration: BoxDecoration(
            color: AppColors.white26,
            borderRadius: BorderRadius.circular(12),
          ),
          padding: const EdgeInsets.all(16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              SizedBox(
                width: 300,
                child: TextFieldWidget(
                  hintText: 'Название',
                  controller: _titleController,
                ),
              ),
              const SizedBox(height: 10),
              const Text(
                'Выбери цвет!',
                style: TextStyle(color: AppColors.white),
              ),
              const SizedBox(height: 16),
              CustomColorPicker(
                pickerColor: currentColor,
                onColorChanged: changeColor,
              ),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  ButtonWidget(
                    title: 'Отменить',
                    color: AppColors.blue,
                    onTap: () => _handleCancel(context),
                  ),
                  const SizedBox(width: 10),
                  ButtonWidget(
                    title: 'Сохранить',
                    color: AppColors.blue,
                    onTap: () => _handleSave(context, ref, labelItem.id),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _handleCancel(BuildContext context) {
    Navigator.of(context).pop();
  }

  void _handleSave(BuildContext context, WidgetRef ref, String labelId) {
    print('title = ${_titleController.text}');
    final newTitle = _titleController.text;
    final newColor = currentColor;

    ref.read(labelNotifierProvider.notifier).updateTitle(labelId, newTitle);
    print('new title = ${labelItem.title}');
    ref.read(labelNotifierProvider.notifier).updateColor(labelId, newColor);

    Navigator.of(context).pop();
  }
}
